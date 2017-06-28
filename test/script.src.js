document.getElementById("#button").addEventListener("click", async() => {
    var conn, result
    conn = maxim.open("ws://localhost/")
    result = await conn.execute("HelloWorld", {
        hello: ", world!"
    })
})