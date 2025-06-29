#!/bin/bash

cargo build --release
BINARY="./target/release/rgrep"


echo "Test 1"
echo "grep J ./cc-test/rockbands.txt"
$BINARY J ./cc-test/rockbands.txt
echo "=========="
echo ""

echo "Test 2"
echo """grep "\d" ./cc-test/test-subdir/BFS1985.txt"""
$BINARY "\d" ./cc-test/test-subdir/BFS1985.txt
echo "=========="
echo ""

echo "Test 3"
echo """grep "\w" ./cc-test/symbols.txt"""
$BINARY "\w" ./cc-test/symbols.txt
echo "=========="
echo ""

echo "Test 4"
echo """grep ^A ./cc-test/rockbands.txt"""
$BINARY ^A ./cc-test/rockbands.txt
echo "=========="
echo ""

echo "Test 5"
echo """grep na$ ./cc-test/rockbands.txt"""
$BINARY na$ rockbands.txt
echo "=========="
echo ""

echo "Test 6"
echo """grep A ./cc-test/rockbands.txt | wc -l"""
$BINARY A ./cc-test/rockbands.txt | wc -l
echo "=========="
echo ""

echo "Test 7"
echo """grep -i A ./cc-test/rockbands.txt | wc -l"""
$BINARY -i A ./cc-test/rockbands.txt | wc -l
echo "=========="
echo ""
