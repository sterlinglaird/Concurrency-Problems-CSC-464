#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>
#include "../Utils/semaphore.h"

unsigned int PASSENGERS_PER_CAR;
unsigned int NUM_PASSENGERS;
unsigned int NUM_ROUNDS;

using std::chrono::system_clock;

semaphore boardQueue;
semaphore unboardQueue;

class Car {
	unsigned int numBoarded = 0;
	std::mutex carMutex;

public:
	void load() {
		for (unsigned int i = 0; i < PASSENGERS_PER_CAR; i++) {
			boardQueue.notify();
		}
	}

	void unload() {
		for (unsigned int i = 0; i < PASSENGERS_PER_CAR; i++) {
			unboardQueue.notify();
		}
	}

	void run() {
		//Maybe condition variable to make more efficent
		while(isEmpty()){}
	}

	void boarded() {
		std::unique_lock<std::mutex> lock{ carMutex };
		numBoarded++;
	}

	void unboarded() {
		std::unique_lock<std::mutex> lock{ carMutex };
		numBoarded--;
	}

	bool isEmpty() {
		std::unique_lock<std::mutex> lock{ carMutex };
		return numBoarded == 0;
	}

};

Car car;

class Passenger {
public:
	void board() {
		boardQueue.wait();
		car.boarded();
	}

	void unboard() {
		unboardQueue.wait();
		car.unboarded();
	}
};

void passengerThread() {
	Passenger passenger;

	passenger.board();
	passenger.unboard();
}

void carThread() {
	for (unsigned int i = 0; i < NUM_ROUNDS; i++) {
		car.load();
		car.run();
		car.unload();
	}
}

int main(int argc, char** args) {
	auto start = system_clock::now();

	PASSENGERS_PER_CAR = atoi(args[1]);
	NUM_PASSENGERS = atoi(args[2]);
	NUM_ROUNDS = atoi(args[3]);

	std::thread car(carThread);

	for (unsigned int i = 0; i < NUM_PASSENGERS; i++) {
		std::thread passenger(passengerThread);
		passenger.detach();
	}

	car.join();

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "passengers per car: " << PASSENGERS_PER_CAR<< std::endl;
	std::cout << "# passengers: " << NUM_PASSENGERS << std::endl;
	std::cout << "# rounds: " << NUM_ROUNDS << std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}