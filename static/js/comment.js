import { DOMError } from "./Error.js";

async function likeComment(data) {
  console.log(data);
  const response = await fetch("/api/like/comment", {
    method: "POST",
    body: JSON.stringify(data),
  });
  const responseData = await response.json();
  if (!response.ok) {
    throw new Error(responseData.message);
  }

  return responseData.message;
}

async function addComment(data) {
  const response = await fetch("/api/add/comment", {
    method: "POST",
    body: JSON.stringify(data),
  });

  const responseData = await response.json();
  if (!response.ok) {
    throw new Error(responseData.message);
  }
  console.log(responseData);

  return responseData.message;
}
//userId, postId, isLike
export function handleLikeComment(commentsId) {
  let commentsEl = document.getElementById(commentsId);
  if (!commentsEl) return;
  commentsEl.addEventListener("click", (e) => {
    const comment = e.target.closest(".comment");
    if (!comment) {
      return;
    }
    const { id } = comment.dataset;
    const data = {
      commentId: +id,
    };
    const likeUp = e.target.closest(".like-up") && 1;
    const likeDown = e.target.closest(".like-down") && -1;

    data.isLike = likeUp || likeDown;
    likeComment(data)
      .catch((err) => console.log(err))
      .then(() => (window.location.href = ""));
  });
}

export function handleCommentForm(formId) {
  const form = document.getElementById(formId);
  if (!form) return;
  const domError = new DOMError(form);
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const { postid } = form.dataset;
    const formData = new FormData(form);
    const data = {
      postId: +postid,
      comment: formData.get("comment"),
    };
    addComment(data)
      .then((data) => {
        domError.writeSucc(data);
        setTimeout(() => {
          window.location.reload();
        }, 800);
      })
      .catch((e) => {
        domError.writeError(e.message);
      });
  });
}
