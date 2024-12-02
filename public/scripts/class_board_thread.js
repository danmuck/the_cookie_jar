let socket = null;

function receiveComment(data) {
    const comment = JSON.parse(data);
    if (comment.Type === "newComment") {
        let editBtn = ""
        const authorId = document.getElementById("username").textContent
        if (document.getElementById("isProf") || (authorId && authorId == comment.AuthorID)) {
            editBtn = '<h6 class="action-btn-edit">Edit</h6>'
        }

        const newCommentHTML = `
        <li class="comment">
            <div class="body">
                <p class="text">` + comment.Text + `</h4>
                <div class="poster">
                    <img src="` + comment.AuthorImageURL  + `" style="width: 60px; height: 60px;">
                    <h6 class="name">` + comment.AuthorID  +`</h6>
                </div>
            </div>
            <div class="actions" data-id="` + comment.ID + `">
                <h6 class="action-btn-like" data-amount="0">Like</h6> ` +
                editBtn +
            `</div>
        </li>
    `;

        document.getElementById('comment-list').insertAdjacentHTML('beforeend', newCommentHTML);
        makeLikeButtonsClickable()
        makeEditButtonsClickable()
    }
    else if (comment.Type === "likeComment") {
        document.querySelectorAll(".comment").forEach(commentElement => {
            const actions = commentElement.querySelector('.actions')
            const id = actions.dataset.id

            if (id === comment.ID) {
                const likeBtn = actions.querySelector('.action-btn-like');
                let likes = likeBtn.dataset.amount;

                if (comment.Liked === "true") {
                    likes++;
                }
                else {
                    likes--;
                }

                likeBtn.dataset.amount = likes;
                if (likes == 0) {
                    likeBtn.textContent = "Like"
                }
                else {
                    likeBtn.textContent = "Like (" + likes + ")"
                }
            }
        });
    }
    else if (comment.Type === "editComment") {
        document.querySelectorAll(".comment").forEach(commentElement => {
            const actions = commentElement.querySelector('.actions')
            const id = actions.dataset.id

            if (id === comment.ID) {
                const textElement = commentElement.querySelector('.text')

                if (textElement.tagName === "P") {
                    textElement.textContent = comment.Text
                }
                else if (textElement.tagName === "TEXTAREA") {
                    textElement.value = comment.Text
                }
            }
        });
    }
}

function sendComment(event) {
    if (socket.readyState === WebSocket.OPEN) {
        const message = {
            Type: "newComment",
            Text: document.getElementById('comment-box').value
        }

        if (message.Text == "") {
            alert("Please add a comment.")
            return;
        }

        document.getElementById('comment-box').value = ""
        socket.send(JSON.stringify(message))
    }
    else {
        console.error("WebSocket is not active.")
    }
}

function likeComment(event) {
    const id = event.currentTarget.closest('.actions').dataset.id
    if (socket.readyState === WebSocket.OPEN) {
        const message = {
            Type: "likeComment",
            ID: id
        }

        socket.send(JSON.stringify(message))
    }
    else {
        console.error("WebSocket is not active.")
    }
}

function editComment(event) {
    const id = event.currentTarget.closest('.actions').dataset.id
    const commentElement = event.currentTarget.closest('.comment')
    const commentText = commentElement.querySelector('.text')

    if (socket.readyState === WebSocket.OPEN) {
        const message = {
            Type: "editComment",
            ID: id,
            Text: commentText.value
        }

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

function editBtnClick(event) {
    const commentElement = event.currentTarget.closest('.comment')
    const commentText = commentElement.querySelector('.text')

    if (event.target.textContent === "Save") {
        event.target.textContent = "Edit";
        editComment(event)

        // Rely on ws to update
        commentText.outerHTML = '<p class="text"></p>';
        return;
    }

    event.target.textContent = "Save";
    commentText.outerHTML = '<textarea class="text" style="width: 100%; height: 100%; margin-top: 14px; resize: vertical;">' + commentText.textContent + '</textarea>';
}

function makeLikeButtonsClickable() {
    document.querySelectorAll(".action-btn-like").forEach(btn => {
        btn.addEventListener('click', likeComment);
    });
}

function makeEditButtonsClickable() {
    document.querySelectorAll(".action-btn-edit").forEach(btn => {
        btn.addEventListener('click', editBtnClick);
    });
}

window.onload = function() {
    socket = new WebSocket(window.location.origin + window.location.pathname + "ws")

    socket.onmessage = function(event) {
        receiveComment(event.data);
    };

    document.getElementById('back-arrow').addEventListener('click', backBtnClick);
    makeLikeButtonsClickable()
    makeEditButtonsClickable()
    document.getElementById('comment-btn').addEventListener('click', sendComment);
}

window.onbeforeunload = function() {
    if (socket.readyState === WebSocket.OPEN) {
        socket.close();
    }
}
