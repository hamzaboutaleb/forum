import { DOMError } from "./Error.js";

async function likeComment(data) {
  const response = await fetch("/api/like/comment", {
    method: "POST",
    body: JSON.stringify(data),
  });
  const responseData = await response.json();
  if (!response.ok) {
    console.log(responseData);
    throw new Error(responseData.message);
  }
  console.log(responseData);

  return responseData.message;
}

async function addComment(data) {
  const response = await fetch("/api/add/comment", {
    method: "POST",
    body: JSON.stringify(data),
  });

  const responseData = await response.json();
  if (!response.ok) {
    console.log(responseData);
    throw new Error(responseData.message);
  }
  console.log(responseData);

  return responseData.message;
}

export function handleLikeComment(comments) {}

export function handleCommentForm(formId) {
  const form = document.getElementById(formId);
  const domError = new DOMError(form);
  if (!form) return;
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const { postid } = form.dataset;
    console.log(form.dataset);
    const formData = new FormData(form);
    const data = {
      postId: +postid,
      comment: formData.get("comment"),
    };
    console.log(data);
    addComment(data)
      .catch((e) => domError.writeError(e))
      .then((data) => {
        domError.writeSucc(data);
        setTimeout(() => {
          window.location.reload();
        }, 800);
      });
  });
}
