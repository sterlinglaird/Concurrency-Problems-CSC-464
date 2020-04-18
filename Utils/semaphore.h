#include <mutex>
#include <condition_variable>

class semaphore
{
private:
	std::mutex mutex;
	std::condition_variable cv;
	unsigned int count = 0;
public:

	semaphore(int count = 0) : count(count){

	}

	void notify() {
		std::lock_guard<std::mutex> lock(mutex);
		count++;
		cv.notify_one();
	}

	void wait() {
		std::unique_lock<std::mutex> lock(mutex);
		while (!count) {
			cv.wait(lock);
		}
		count--;
	}

	bool try_wait() {
		std::lock_guard<std::mutex> lock(mutex);
		if (count) {
			count--;
			return true;
		}
		return false;
	}
};
