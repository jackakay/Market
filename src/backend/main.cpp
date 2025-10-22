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
	
	svr.Get("/api/getaskprice", [askQueue](const httplib::Request&, httplib::Response& res) {
      		int cheapestAsk = askQueue.top().price;
		json data = { {"price", (cheapestAsk)}	};
        	res.set_content(data.dump(), "application/json");
    	});

	svr.Get("/api/getbidprice", [bidQueue](const httplib::Request&, httplib::Response& res){
		int highestBid = bidQueue.top().price;		
		json data = {{"price", highestBid}};
		res.set_content(data.dump(), "application/json");
	});


	svr.listen("0.0.0.0", 8080);
}







