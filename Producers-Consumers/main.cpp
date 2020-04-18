#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

typedef int event;

//Does not work when buffer size is 1 because of the way I defined "empty" and "full"
 unsigned int BUFFER_SIZE = 100;
 unsigned int NUM_PRODUCERS = 2;
 unsigned int NUM_CONSUMERS = 10;
 unsigned int NUM_PRODUCED = 10000; //each

 using std::chrono::system_clock;

class Buffer {

	std::unique_ptr<event> buffer;
	int startIdx = 0;
	int openIdx = 0;

	std::mutex mutex;
	std::condition_variable condProducer;
	std::condition_variable condConsumer;
	
public:
	Buffer(unsigned int size) {
		buffer = std::unique_ptr<event>(new event[size]);
	}
	event getEvent() {
		std::unique_lock<std::mutex> lock{ mutex };

		while (isEmpty()) {
			condConsumer.wait(lock);
		}

		event ev = buffer.get()[startIdx];

		startIdx = (startIdx + 1) % BUFFER_SIZE;
		condProducer.notify_one();
		return ev;
	}

	void addEvent(event ev) {
		std::unique_lock<std::mutex> lock{ mutex };

		while (isFull()){
			condProducer.wait(lock);
		}

		buffer.get()[openIdx] = ev;
		openIdx = (openIdx + 1) % BUFFER_SIZE;

		condConsumer.notify_one();
	}

	bool isEmpty() {
		return startIdx == openIdx;
	}

	bool isFull() {
		return startIdx == (openIdx + 1) % BUFFER_SIZE;
	}
};

std::unique_ptr<Buffer> buffer;

event waitForEvent() {
	int randomNum = std::rand();
	return event(randomNum % 10000);
}

void consumeEvent(event ev) {

}

void producer(int id) {
	for (unsigned int i = 0; i < NUM_PRODUCED; i++) {
		event ev = waitForEvent();

		buffer->addEvent(ev);
	}
}

void consumer(int id) {
	while (true) {
		event ev = buffer->getEvent();
		consumeEvent(ev);
	}
}

int main(int argc, char** args) {
	auto start = system_clock::now();

	BUFFER_SIZE = atoi(args[1]);
	NUM_PRODUCERS = atoi(args[2]);
	NUM_CONSUMERS = atoi(args[3]);
	NUM_PRODUCED = atoi(args[4]);

	buffer = std::unique_ptr<Buffer>(new Buffer(BUFFER_SIZE));

	auto producers = std::unique_ptr<std::thread[]>(new std::thread[NUM_PRODUCERS]);
	
	for (unsigned int i = 0; i < NUM_PRODUCERS; i++) {
		producers[i] = std::thread(producer, i);
	}

	for (unsigned int i = 0; i < NUM_CONSUMERS; i++) {
		std::thread consumer(consumer, i);
		consumer.detach();
	}

	for (unsigned int i = 0; i < NUM_PRODUCERS; i++) {
		producers[i].join();
	}

	auto end = system_clock::now();
	auto difference = end - start;
	auto milliseconds = std::chrono::duration_cast<std::chrono::milliseconds>(difference).count();

	std::cout << "buffer size: " << BUFFER_SIZE << std::endl;
	std::cout << "# producers: " << NUM_PRODUCERS << std::endl;
	std::cout << "# consumers: " << NUM_CONSUMERS << std::endl;
	std::cout << "# produced per producer: " << NUM_PRODUCED << std::endl;
	std::cout << "Time taken: " << milliseconds << "ms" << std::endl;

	return 0;
}