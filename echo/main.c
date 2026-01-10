#include <asm-generic/errno.h>
#include <errno.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <poll.h>
#include <stdbool.h>

int udp_echo_server();

void handle(const char *reason) {
	printf("Error No: %d\n",errno);
	printf("Failed: %s\n", reason);
}

int main(int argc, char *argv[]) {
	int sock_fd , nfds;
	struct sockaddr_in my_addr;
	struct sockaddr_in client_addr;
	socklen_t client_addr_len = sizeof(client_addr); 
	struct pollfd fds[200];

	// UDP
	if(argc > 1 && strcmp(argv[1], "-udp") == 0){
		return udp_echo_server();
	}	

	// TCP
	sock_fd =  socket(AF_INET,SOCK_STREAM,0);
	if(sock_fd == -1){
		return -1;
	}

	memset(&my_addr, 0, sizeof(struct sockaddr_in));
	my_addr.sin_family = AF_INET;
	my_addr.sin_addr.s_addr = INADDR_ANY;
	my_addr.sin_port = htons(8000);

	if (bind(sock_fd, (struct sockaddr*)&my_addr, sizeof(my_addr)) != 0){
		handle("BIND");
		return -1;
	}

	if(listen(sock_fd,32) != 0){
		handle("LISTEN");
		return -1;
	}

	fds[0].fd = sock_fd;
	fds[0].events = POLLIN;
	nfds = 1;

	printf("TCP Echo server listening on :8000\n");
	int timeout = 30 * 1000;
	do {
  		bool COMPRESS_ARRAY = false;
		int ret = poll(fds, nfds, timeout);
		if (ret < 0) {
			printf("poll() failed... Exiting\n");
			break;
		}

		if (ret == 0){
			printf("Timeout. Exiting...\n");
			break;
		}

		int current_len = nfds;
		for(int i=0;i<current_len;i++){
			if(fds[i].revents == POLLIN) {
				if ( fds[i].fd == sock_fd ) {
					int client_fd = accept(sock_fd,(struct sockaddr*)&client_addr, &client_addr_len);
					if(client_fd == -1){
						handle("ACCEPT");
						continue;
					}

					printf("Accepted connection from %s\n", inet_ntoa(client_addr.sin_addr));
					fds[nfds].fd = client_fd;
					fds[nfds].events = POLLIN;
					nfds++;
				} else {
					char buf[500];
					size_t size = recv(fds[i].fd, buf, 500,0);
					if(size == 0){
						fds[i].fd = -1;
						COMPRESS_ARRAY = true;
						continue;
					}
					buf[size] = '\0';

					if(send(fds[i].fd, buf, size, 0) == -1){
						printf("Error while sending data: %d\n", errno);
						fds[i].fd = -1;
						COMPRESS_ARRAY = true;
					}
				}
			}
		}

		if (COMPRESS_ARRAY == true) {
			for(int i=0;i<nfds;i++){
				if(fds[i].fd == -1){
					printf("Client closed\n");
					for(int j=i;j<nfds-1;j++){
						fds[j] = fds[j+1];
					}
					i--;
					nfds--;
				}
			}
		}
	} while(1);



	for(int i=0;i<nfds;i++){
		close(fds[i].fd);
	}
	close(sock_fd);


	return 0;
}

int udp_echo_server(){
	int sock_fd , nfds;
	struct sockaddr_in my_addr;
	struct sockaddr_in client_addr;
	socklen_t client_addr_len = sizeof(client_addr); 
	struct pollfd fds[200];

	sock_fd =  socket(AF_INET,SOCK_DGRAM,0);
	if(sock_fd == -1){
		return -1;
	}

	memset(&my_addr, 0, sizeof(struct sockaddr_in));
	my_addr.sin_family = AF_INET;
	my_addr.sin_addr.s_addr = INADDR_ANY;
	my_addr.sin_port = htons(8000);

	if (bind(sock_fd, (struct sockaddr*)&my_addr, sizeof(struct sockaddr_in)) != 0){
		handle("BIND");
		return -1;
	}

	printf("UDP Echo server listening on :8000\n");

	while(1){
		char buf[500];
		size_t len = recvfrom(sock_fd, buf , 500 ,0 ,(struct sockaddr*)&client_addr, &client_addr_len);

		printf("Received : %s\n", buf);

		sendto(sock_fd, buf , len ,0 , (struct sockaddr*)&client_addr, client_addr_len);
	}

	close(sock_fd);

	return 0;
}
