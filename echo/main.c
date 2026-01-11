#include <asm-generic/errno.h>
#include <errno.h>
#include <signal.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/poll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <poll.h>
#include <stdbool.h>

volatile sig_atomic_t shutdown_server = 0;
const int TIMEOUT = 3 * 60 * 1000;

int udp_echo_server();
int tcp_echo_server();

void handle(const char *reason) {
	printf("Error No: %d\n",errno);
	printf("Failed: %s\n", reason);
}

void shutdown_handler(int signum) {
	if (signum == SIGTERM || signum == SIGINT) {
		printf("shuting down...\n");
		shutdown_server = 1;
		errno = 0;
	}
}

int main(int argc, char *argv[]) {
	signal(SIGTERM,shutdown_handler);
	signal(SIGINT, shutdown_handler);
	// UDP
	if(argc > 1 && strcmp(argv[1], "-udp") == 0){
		return udp_echo_server();
	}

	// TCP
	return tcp_echo_server();
}

int tcp_echo_server(){
	int sock_fd , nfds;
	struct sockaddr_in my_addr;
	struct sockaddr_in client_addr;
	socklen_t client_addr_len = sizeof(client_addr); 
	struct pollfd fds[200];

	sock_fd =  socket(AF_INET,SOCK_STREAM,0);
	if(sock_fd == -1){
		return errno;
	}

	memset(&my_addr, 0, sizeof(struct sockaddr_in));
	my_addr.sin_family = AF_INET;
	my_addr.sin_addr.s_addr = INADDR_ANY;
	my_addr.sin_port = htons(8000);

	if (bind(sock_fd, (struct sockaddr*)&my_addr, sizeof(my_addr)) != 0){
		handle("BIND");
		return errno;
	}

	if(listen(sock_fd,32) != 0){
		handle("LISTEN");
		return errno;
	}

	fds[0].fd = sock_fd;
	fds[0].events = POLLIN;
	nfds = 1;

	printf("TCP Echo server listening on :8000\n");
	do {
  		bool COMPRESS_ARRAY = false;
		int ret = poll(fds, nfds, TIMEOUT);
		if (ret < 0) {
			if (shutdown_server != 1){
				printf("poll() failed... Exiting\n");
			}
			break;
		}

		if (ret == 0){
			printf("Timeout. Exiting...\n");
			break;
		}

		int current_len = nfds;
		for(int i=0;i<current_len;i++){
			if(fds[i].revents & POLLIN){
				if (fds[i].fd == sock_fd){
					int client_fd = accept(sock_fd,(struct sockaddr*)&client_addr, &client_addr_len);
					if(client_fd == -1){
						handle("ACCEPT");
						continue;
					}

					printf("Accepted connection from %s:%d\n", inet_ntoa(client_addr.sin_addr), client_addr.sin_port);
					fds[nfds].fd = client_fd;
					fds[nfds].events = POLLIN;
					nfds++;
				} else {
					int capacity = 1024;
					char *data = malloc(capacity);
					size_t total_size = 0;
					size_t size = 0; 
					char buf[1024];
					do{
						memset(buf,0, 1024);
						size = recv(fds[i].fd, buf, 1024,0);
						total_size += size;

						if(total_size >= capacity){
							capacity = capacity * 2 + 1;
							char *new  = realloc(data, capacity);
							if (new == NULL){
								handle("Error allocating new memory");
								total_size = 0;
								free(new);
								break;
							}
							data = new;
						}
						strcat(data, buf);
						data[total_size] = '\0';
					}while(size == 1024);

					if(total_size <= 0){
						if (total_size < 0 ){
							handle("Error while receiving data");
						}
						fds[i].fd = -1;
						COMPRESS_ARRAY = true;
						free(data);
						break;
					}

					data[total_size] = '\0';
					if(send(fds[i].fd, data, total_size, 0) == -1){
						handle("Error while sending data");
						fds[i].fd = -1;
						COMPRESS_ARRAY = true;
					}
					free(data);
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
	} while(shutdown_server == 0);

	for(int i=0;i<nfds;i++){
		close(fds[i].fd);
	}
	close(sock_fd);

	return errno;
}

int udp_echo_server(){
	int sock_fd;
	struct sockaddr_in my_addr;
	struct sockaddr_in client_addr;
	socklen_t client_addr_len = sizeof(client_addr); 
	struct pollfd fd[1];

	sock_fd = socket(AF_INET,SOCK_DGRAM,0);
	if(sock_fd == -1){
		return errno;
	}

	memset(&my_addr, 0, sizeof(struct sockaddr_in));
	my_addr.sin_family = AF_INET;
	my_addr.sin_addr.s_addr = INADDR_ANY;
	my_addr.sin_port = htons(8000);

	if (bind(sock_fd, (struct sockaddr*)&my_addr, sizeof(struct sockaddr_in)) != 0){
		handle("BIND");
		return errno;
	}

	fd[0].fd = sock_fd;
	fd[0].events = POLLIN;

	printf("UDP Echo server listening on :8000\n");
	do{
		int ret = poll(fd, 1, TIMEOUT);
		if (ret < 0) {
			if (shutdown_server != 1){
				printf("poll() failed... Exiting\n");
			}
			break;
		}

		if (ret == 0){
			printf("Timeout. Exiting...\n");
			break;
		}

		if(fd[0].revents & POLLIN){
			char data[65536];
			size_t size = recvfrom(sock_fd, data, 65536 ,0 ,(struct sockaddr*)&client_addr, &client_addr_len);
			printf("Received message from %s:%d\n", inet_ntoa(client_addr.sin_addr), client_addr.sin_port);
			data[size] = '\0';
			if(sendto(sock_fd, data , size ,0 , (struct sockaddr*)&client_addr, client_addr_len)== -1){
				handle("Error while sending data");
			}
		}
	}while(shutdown_server == 0);

	close(sock_fd);

	return 0;
}
