let socket = null;
let state = "waiting";
let answered = false;

function receiveMessage(data) {
    const message = JSON.parse(data)
    switch (message.Type) {
        case "playerList":
            const playerListElement = document.getElementById("player-list")

            let playerListItems = ""
            for (let player in message.Players) {
                if (message.Players.hasOwnProperty(player)) {
                    const state = message.Players[player]
                    if (state == 0) {
                        playerListItems += '<li style="color: red;">'
                    }
                    else {
                        playerListItems += '<li style="color: green;">'
                    }
                    playerListItems += player + '</li>\n'
                }
            }

            playerListElement.innerHTML = playerListItems;
            break;

        case "error":
            alert(message.Text)
            break;

        case "startCountdown":
            hideAllContentButOne("start-countdown")
            const startTimer = document.getElementById("start-timer");
            startTimer.innerText = message.Seconds;
            break;

        case "questionPrep":
            hideAllContentButOne("question-prep")
            const prepQuestionText = document.getElementById("question-prep-text");
            prepQuestionText.innerText = message.Text;
            state = "prep";
            break;

        case "question":
            hideAllContentButOne("question")
            const questionText = document.getElementById("question-text");
            questionText.innerText = message.Text;
            state = "question";
            answered = false;

            const optionListElement = document.getElementById("question-answers")
            let optionListItems = ""
            let count = 0;
            for (let option in message.Options) {
                if (message.Options.hasOwnProperty(option)) {
                    const optionText = message.Options[option]
                    optionListItems += '<button id="option-' + count + '">' + optionText + '</button>';
                    count++;
                }
            }
            optionListElement.innerHTML = optionListItems;
            makeOptionsClickable();
            break;

        case "leaderboard":
            hideAllContentButOne("leaderboard")
            const topThreeElement = document.getElementById("top-three")

            const topThree = Object.entries(message.Scores).sort((a, b) => b[1] - a[1]).slice(0, 3);
            let topThreeImages = ""
            for (let player in topThree) {
                if (topThree.hasOwnProperty(player)) {
                    topThreeImages += `
                        <div class="user">
                            <img src="/pfp/` + topThree[player][0] + `" style="width: 60px; height: 60px;">
                            <label>` + topThree[player][0] + `</label>
                            <label>|</label>
                            <label>` + topThree[player][1] + ` Points</label>
                        </div>
                    `
                }
            }
            topThreeElement.innerHTML = topThreeImages;
            break;

        case "score":
            const pointsElement = document.getElementById("points")
            pointsElement.textContent = "You Have " + message.Score + " Points"
            const correctOptionElement = document.getElementById("correct-ans")
            correctOptionElement.textContent = "Correct Option was Option " + message.CorrectOption
            break;

        case "questionCountdown":
            if (state === "prep") {
                const prepTimer = document.getElementById("question-prep-timer");
                prepTimer.innerText = message.Seconds;
            }
            else if (state === "question") {
                const questionTimer = document.getElementById("question-timer");
                questionTimer.innerText = message.Seconds;
            }
            break;
    }
}

function sendMessage(data) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(data)
    }
    else {
        console.error("WebSocket is not active.")
    }
}

function hideAllContentButOne(id) {
    const contentItem = document.getElementById("content")
    const children = contentItem.children

    for (let i = 0; i < children.length; i++) {
        if (children[i].id !== id) {
            children[i].style.display = 'none';
        }
        else {
            children[i].style.display = 'block';
        }
    }
}

function makeOptionsClickable() {
    const children = document.getElementById("question-answers").children

    for (let i = 0; i < children.length; i++) {
        const buttonElement = children[i]

        buttonElement.onclick = function(event) {
            if (answered) {
                return;
            }

            const optionNumber = event.target.id.charAt(event.target.id.length - 1)
            const message = {
                Type: "questionAnswer",
                Option: optionNumber

            }
            socket.send(JSON.stringify(message));
            
            event.target.style.backgroundColor = "green";
            answered = true;
        }
    }
}

function startBtnClick(event) {
    const message = {
        Type: "start"
    }
    socket.send(JSON.stringify(message))
}

function readyUpBtnClick(event) {
    const message = {
        Type: "ready"
    }
    socket.send(JSON.stringify(message))
}

window.onload = function() {
    socket = new WebSocket(window.location.origin + window.location.pathname + "ws")

    socket.onopen = function() {
        const message = {
            Type: "join"
        }
        socket.send(JSON.stringify(message))
    }

    socket.onmessage = function(event) {
        receiveMessage(event.data);
    };

    if (document.getElementById('start-btn')) {
        document.getElementById('start-btn').addEventListener('click', startBtnClick);
    }
    document.getElementById('ready-up-btn').addEventListener('click', readyUpBtnClick);
}

window.onbeforeunload = function() {
    if (socket.readyState === WebSocket.OPEN) {
        socket.close();
    }
}
