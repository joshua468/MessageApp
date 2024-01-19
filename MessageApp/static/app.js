
document.addEventListener("DOMContentLoaded", function () {
    loadMessages();
});

function loadMessages() {
    fetch("/messages")
        .then(response => response.json())
        .then(messages => {
            const messageContainer = document.getElementById("message-container");
            messageContainer.innerHTML = "";
            messages.forEach(message => {
                const messageElement = document.createElement("div");
                messageElement.classList.add("message");
                messageElement.innerHTML = `<strong>${message.username}</strong>: ${message.content}`;
                messageContainer.appendChild(messageElement);
            });
            messageContainer.scrollTop = messageContainer.scrollHeight;
        })
        .catch(error => console.error("Error fetching messages:", error));
}

function sendMessage() {
    const username = document.getElementById("username").value;
    const content = document.getElementById("content").value;

    if (username && content) {
        const message = {
            username: username,
            content: content
        };

        fetch("/messages", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(message),
        })
            .then(() => {
                loadMessages();
                document.getElementById("content").value = "";
            })
            .catch(error => console.error("Error sending message:", error));
    } else {
        alert("Username and content cannot be empty.");
    }
}
