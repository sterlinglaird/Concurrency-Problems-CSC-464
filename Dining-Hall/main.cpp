#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>
#include "../Utils/semaphore.h"

using std::chrono::system_clock;

unsigned int NUM_STUDENTS = 1000;
unsigned int MAX_WAIT_TIME = 50;
unsigned int NUM_ROUNDS = 50;

semaphore okToLeave(0);

std::mutex mutex;

unsigned int numEating;
unsigned int numReadyToLeave;

void dine(unsigned int i) {
	std::unique_lock<std::mutex> lock{ mutex };
	numEating++;
	if (numEating == 2 && numReadyToLeave == 1) {
		okToLeave.notify();
	}

	lock.unlock();

	//Eat!!
	
	lock.lock();

	numEating--; //Done eating
}

void leave(unsigned int i) {
	std::unique_lock<std::mutex> lock{ mutex };
	numReadyToLeave++;

	if (numEating == 1 && numReadyToLeave == 1) {
		lock.unlock();
		okToLeave.wait();
		lock.lock();
	}
	// If there is one other student waiting to leave
	else if (numEating == 0 && numReadyToLeave == 2) {
		okToLeave.notify();
		numReadyToLeave -= 2;
	}
	else {
		numReadyToLeave--;
	}
}

void studentThread(unsigned int i) {
	for (unsigned int i = 0; i < NUM_ROUNDS; i++) {
		dine(i);

		leave(i);
	}
}

int main(int argc, char** args) {
	auto start = system_clock::now();

	NUM_STUDENTS = atoi(args[1]);
	NUM_ROUNDS = atoi(args[2]);

	auto students = std::unique_ptr<std::thread[]>(new std::thread[NUM_STUDENTS]);

	for (unsigned int i = 0; i < NUM_STUDENTS; i++) {
		students[i] = std::thread(studentThread, i);
	}

	for (unsigned int i = 0; i < NUM_STUDENTS; i++) {
		students[i].join();
	}

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "# students: " << NUM_STUDENTS << std::endl;
	std::cout << "# rounds: " << NUM_ROUNDS<< std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}