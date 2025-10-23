setInterval(async () => {
  try {
    const response = await fetch("/api/price"); // or "http://localhost:8080/price" if different origin
    if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);

    const data = await response.json();
    console.log("Latest price:", data.price);

    // Example: update the page
    document.getElementById("price").textContent = data.price.toFixed(2);
    document.getElementById("askprice").textContent = data.ask.toFixed(2);
    document.getElementById("bidprice").textContent = data.bid.toFixed(2);
  } catch (err) {
    console.error("Failed to fetch price:", err);
  }
}, 2000);