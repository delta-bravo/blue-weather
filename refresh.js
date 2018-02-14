function refresh() {
    console.log("Refreshing temperature value");
    $.ajax({
        url: '/temperature',
        method: 'GET',
        success: function (result) {
            document.getElementById("temp").innerText = "Ambient temperature is " + result + " degrees celsius"
        },
        error: function (xhr, status, error) {
            console.log(error);
        }
    });
}

setInterval(refresh, 5000);
