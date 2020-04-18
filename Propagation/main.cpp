#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>
#include "../Utils/semaphore.h"

using std::chrono::system_clock;

unsigned int NUM_COPIES = 100;
unsigned int MAX_DELAY = 50;
unsigned int NUM_ROUNDS = 10;

class Resource {
	int value = 0;

	std::mutex mutex;
	semaphore updatedValue;

public:
	void setValue(int newValue) {
		std::unique_lock<std::mutex> lock{ mutex };
		value = newValue;
		updatedValue.notify();
	}

	int waitForNewValue() {
		updatedValue.wait();
		std::unique_lock<std::mutex> lock{ mutex };
		int retVal = value;
		return retVal;
	}
};

std::unique_ptr<Resource[]> resources;

void copyThread(unsigned int id, bool original) {
	int value = 0; //The local copy
	while (true) {
		value = resources[id].waitForNewValue();

		if (original) {
			for (unsigned int propagateId = 0; propagateId < NUM_COPIES; propagateId++) {
				if (propagateId == id) {
					continue;
				}

				resources[propagateId].setValue(value);
			}
		}
	}

}
std::mutex startMut;
void mutatorThread(unsigned int id) {
	std::unique_lock<std::mutex> lock{ startMut };
	for (unsigned int i = 0; i < NUM_ROUNDS; i++) {
		if (MAX_DELAY != 0) {
			int msToSleep = std::rand() % MAX_DELAY;
			std::this_thread::sleep_for(std::chrono::milliseconds(msToSleep));
		}

		int newValue = std::rand() % 10000; //Makes number more manageable to look at

		resources[id].setValue(newValue);
	}
}

int main(int argc, char** args) {
	auto start = system_clock::now();

	NUM_COPIES = atoi(args[1]);
	MAX_DELAY = atoi(args[2]);
	NUM_ROUNDS = atoi(args[3]);

	resources = std::unique_ptr<Resource[]>(new Resource[NUM_COPIES]);
	
	//Dont start mutating until all threads are made
	startMut.lock();

	//Only mutate one copy (the original)
	auto mutator = std::thread(mutatorThread, 0);

	for (unsigned int i = 0; i < NUM_COPIES; i++) {
		auto copy = std::thread(copyThread, i, i==0);
		copy.detach();
	}

	startMut.unlock();

	mutator.join();

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "# students: " << NUM_COPIES << std::endl;
	std::cout << "max delay: " << MAX_DELAY << std::endl;
	std::cout << "# rounds: " << NUM_ROUNDS << std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}