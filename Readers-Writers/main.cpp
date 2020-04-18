#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>
#include "../Utils/semaphore.h"

//No starve readers-writers

using std::chrono::system_clock;

unsigned int NUM_READERS = 10;
unsigned int NUM_WRITERS = 10;
unsigned int NUM_ACTIONS = 100; //Per reader/writer

class Resource {
	int value = 0;
	unsigned int numReaders = 0;
	std::mutex mutex;
	semaphore empty;

public:
	Resource() : empty(1) {

	}

	int read() {
		std::unique_lock<std::mutex> lock{ mutex };
		numReaders++;

		//first one takes the empty flag
		if (numReaders == 1) {
			empty.wait();
		}
		lock.unlock();

		int ret = value;

		lock.lock();
		numReaders--;

		//Last one sends empty signal
		if (numReaders == 0)
		{
			empty.notify();
		}

		return ret;
	}

	void write(int value) {
		empty.wait();
		this->value = value;
		empty.notify();
	}
};

Resource resource;

std::mutex coutMut;

void readerThread(unsigned int idx) {
	for (unsigned int i = 0; i < NUM_ACTIONS; i++) {
		int value = resource.read();
	}
}

void writerThread(unsigned int idx) {
	std::srand(idx);
	for (unsigned int i = 0; i < NUM_ACTIONS; i++) {
		int randomNum = std::rand();
		resource.write(randomNum);
	}
}

int main(int argc, char** args) {
	auto start = system_clock::now();

	NUM_READERS = atoi(args[1]);
	NUM_WRITERS = atoi(args[2]);
	NUM_ACTIONS = atoi(args[3]);

	auto readers = std::unique_ptr<std::thread[]>(new std::thread[NUM_READERS]);
	for (unsigned int i = 0; i < NUM_READERS; i++) {
		readers[i] = std::thread(readerThread, i);
	}

	auto writers = std::unique_ptr<std::thread[]>(new std::thread[NUM_WRITERS]);
	for (unsigned int i = 0; i < NUM_WRITERS; i++) {
		writers[i] = std::thread(writerThread, i);
	}

	//Wait for threads to finish
	for (unsigned int i = 0; i < NUM_WRITERS; i++) {
		writers[i].join();
	}
	for (unsigned int i = 0; i < NUM_READERS; i++) {
		readers[i].join();
	}

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "# readers: " << NUM_READERS << std::endl;
	std::cout << "# writers: " << NUM_WRITERS << std::endl;
	std::cout << "# actions: " << NUM_ACTIONS << std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}