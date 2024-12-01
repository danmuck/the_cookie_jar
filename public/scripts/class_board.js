let socket = null;

function receiveThread(data) {
    const thread = JSON.parse(data);
    if (thread.Type != "newThread") {
        return;
    }

    const newThreadHTML = `
        <li class='thread'>
            <div class='title-date'>
                <h4 class='title' data-id='` + thread.ID + `'>` + thread.Title + `</h4>
                <h5 class='date'>` + thread.Date + `</h5>
            </div>
            <div class='poster'>
                <img src='` + thread.AuthorImageURL + `' style='width: 60px; height: 60px;'>
                <h6 class='name'>` + thread.AuthorID + `</h6>
            </div>
        </li>
    `;

    document.getElementById('normal-thread-list').insertAdjacentHTML('beforeend', newThreadHTML);
    makeThreadTitlesClickable();
}

function sendThread(event) {
    if (socket.readyState === WebSocket.OPEN) {
        const message = {
            Type: "newThread",
            Title: document.getElementById('comment-title-box').value,
            Comment: document.getElementById('comment-box').value
        }

        if (message.Title == "") {
            alert("Please add a thread title.")
            return;
        }
        if (message.Comment == "") {
            alert("Please add a thread comment.")
            return;
        }

        document.getElementById('comment-title-box').value = ""
        document.getElementById('comment-box').value = ""
        socket.send(JSON.stringify(message))
    }
    else {
        console.error("WebSocket is not active.")
    }
}

function backBtnClick(event) {
    let urlParts = window.location.pathname.split('/');
    if (urlParts[urlParts.length - 1] === "") {
        urlParts.pop();
    }
    urlParts.pop();

    window.location.assign(urlParts.join('/'));
}

function threadClick(event) {
    const threadItem = event.currentTarget;
    const id = threadItem.dataset.id

    window.location.assign(window.location.origin + window.location.pathname + id)
}

function makeThreadTitlesClickable() {
    document.querySelectorAll(".thread .title-date .title").forEach(threadItem => {
        threadItem.addEventListener('click', threadClick);
    });
}

window.onload = function() {
    socket = new WebSocket(window.location.origin + window.location.pathname + "ws")

    socket.onmessage = function(event) {
        receiveThread(event.data);
    };

    document.getElementById('back-arrow').addEventListener('click', backBtnClick);
    document.getElementById('post-btn').addEventListener('click', sendThread);
    makeThreadTitlesClickable();
}

window.onbeforeunload = function() {
    if (socket.readyState === WebSocket.OPEN) {
        socket.close();
    }
}
