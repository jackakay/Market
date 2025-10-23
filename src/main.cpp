#include "Utils/httplib.h"
#include "Utils/json.hpp"
#include <set>
#include <string>
#include <chrono>
#include <stdexcept>
#include <random>


using json = nlohmann::json;
using namespace std; 

std::random_device rd;
std::mt19937 gen;

std::uniform_int_distribution<> disInt(95,105);
struct Order{
	double price;
	int shares;
	string userID;
	long long timestamp;
};

struct AscendingOrder{

	bool operator()(const Order& a, const Order& b) const{
		if(a.price != b.price) return a.price < b.price;
		else return a.timestamp < b.timestamp;
	}
};

struct DescendingOrder{
	bool operator()(const Order& a, const Order& b) const{
		if(a.price != b.price) return a.price > b.price;
		else return a.timestamp < b.timestamp;
	}
};


template <typename Comparator>
class Queue {
private:
    std::multiset<Order, Comparator> orders;

    long long currentTimestamp() {
        return std::chrono::duration_cast<std::chrono::milliseconds>(
            std::chrono::system_clock::now().time_since_epoch()
        ).count();
    }

public:
    // Add a new order
    void push(double price, int shares, const std::string& userID) {
        orders.insert({price, shares, userID, currentTimestamp()});
    }

    // Get the best (lowest) ask
    Order top() const {
        if (orders.empty()) throw std::runtime_error("Queue is empty");
        return *orders.begin();
    }
 
    // Remove the best ask (fully executed)
    void pop() {
        if (!orders.empty())
            orders.erase(orders.begin());
    }

    // Check if empty
    bool empty() const {
        return orders.empty();
    }
};

int main() {
	httplib::Server svr;
	Queue<DescendingOrder> askQueue;
	Queue<AscendingOrder> bidQueue;
	//Here this is an ascending order queue, so when trying to fill orders it will work from the cheapest first.*
	for(int i = 0; i <10; i++){
		askQueue.push(disInt(gen), 5, "user");
	}
	for(int i = 0; i <10; i++){
		bidQueue.push(disInt(gen), 5, "user");
	}
	

	svr.Options(".*", [](const httplib::Request &, httplib::Response &res) {
        res.set_header("Access-Control-Allow-Origin", "*");
        res.set_header("Access-Control-Allow-Headers", "*");
        res.set_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
        res.status = 204; // No Content
    });

	svr.set_mount_point("/", "frontend");
	svr.Get("/api/price", [askQueue, bidQueue](const httplib::Request&, httplib::Response& res) {

		res.set_header("Content-Type", "text/event-stream");
        res.set_header("Cache-Control", "no-cache");
        res.set_header("Connection", "keep-alive");

      	
		double cheapestAsk = askQueue.top().price;
		double highestBid = bidQueue.top().price;	
		double price = (cheapestAsk+highestBid)/2;

		
		json data = { 
			{"price", (price)},
			{"ask", (cheapestAsk)},
			{"bid", (highestBid)}
		};
		res.set_content(data.dump(), "application/json");
    	});

	


	svr.listen("0.0.0.0", 8080);
}







