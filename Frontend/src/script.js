const buyButton = document.getElementById('buyButton');
const priceInput = document.getElementById('priceInput');
const amountInput = document.getElementById('amount');
const price = document.getElementById('price');

const myHeaders = new Headers();
myHeaders.append("Content-Type", "application/json");

const requestOptionsPrice = {
  method: "GET",
  redirect: "follow"
};


buyButton.addEventListener('click', () => {
        // TODO: Add buy logic here'
        const raw = JSON.stringify({
            "price": Number(priceInput.value),
            "amount": Number(amountInput.value),
            "stock": "JD",
            "type": "MARKET_BUY"
        });
        const requestOptions = {
  method: "POST",
  headers: myHeaders,
  body: raw,
  redirect: "follow"
};
        console.log('Buy button clicked');
        console.log(requestOptions);
        fetch("http://localhost:8080/buy", requestOptions)
        
  .then((response) => response.text())
  .then((result) => buyProcess(result))
  .catch((error) => console.error(error));
});
function setPrice(newPrice) {
    price.innerText = `Current Price: $${newPrice}`;
}
function buyProcess(json){
    const buyData = JSON.parse(json);
    let averagePrice;
    let totalShares = 0;
    let tempAmount = 0;
    let count = 0;
    for (const item of buyData) {
        count += 1;
        tempAmount += (item.price * item.amount);
        totalShares += item.amount;
    }
    averagePrice = (tempAmount / totalShares).toFixed(2);
    setPrice(averagePrice);
    alert(`Bought ${totalShares} shares at average price: $${averagePrice}`);
}