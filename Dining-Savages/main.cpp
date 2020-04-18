#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

using std::chrono::system_clock;

unsigned int NUM_SAVAGES = 10;
unsigned int MAX_SERVINGS = 5;
unsigned int NUM_COOK_ITERS = 10;

enum Serving {food};

class Pot {
	unsigned int numServings = 0;

	std::mutex mutex;
	std::condition_variable fullPot;
	std::condition_variable emptyPot;

public:
	void waitUntilEmpty() {
		std::unique_lock<std::mutex> lock{ mutex };
		while (numServings > 0) {
			emptyPot.wait(lock);
		}
	}

	void fillPot(unsigned int numServings){
		std::unique_lock<std::mutex> lock{ mutex };
		this->numServings = numServings;
		fullPot.notify_all();
	}

	Serving getServing(int idx) {
		std::unique_lock<std::mutex> lock{ mutex };
		while (numServings <= 0) {
			emptyPot.notify_one();
			fullPot.wait(lock);
		}

		numServings--;
		return Serving::food;
	}
};

Pot pot;

void eat(Serving serving) {

}

void cook() {
	for (unsigned int i = 0; i < NUM_COOK_ITERS; i++) {
		pot.waitUntilEmpty();
		pot.fillPot(MAX_SERVINGS);
	}
}

void savage(unsigned int idx) {
	while (true) {
		Serving serving = pot.getServing(idx);
		eat(serving);
	}

}

int main(int argc, char** args) {
	auto start = system_clock::now();

	NUM_SAVAGES = atoi(args[1]);
	MAX_SERVINGS = atoi(args[2]);
	NUM_COOK_ITERS = atoi(args[3]);

	for (unsigned int i = 0; i < NUM_SAVAGES; i++) {
		std::thread savage(savage, i);
		savage.detach();
	}

	std::thread cook(cook);

	cook.join();

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "# savages: " << NUM_SAVAGES << std::endl;
	std::cout << "Max servings: " << MAX_SERVINGS << std::endl;
	std::cout << "# cook iterations: " << NUM_COOK_ITERS << std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}