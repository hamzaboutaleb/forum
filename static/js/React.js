async function reactToPost(postId, isLike, countEl) {
  try {
    const response = await fetch(`/api/react`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        postId,
        isLike,
      }),
    });

    const result = await response.json();
    if (response.ok) {
      console.log(result.message);
    } else {
      console.error(`Error: ${response.status} - ${result.message}`);
    }
    countEffect(countEl, response.ok);
    return result;
  } catch (error) {
    console.log(error);
  }
}

export function ReactHandler() {
  let posts = document.querySelectorAll(".post");

  posts.forEach((post) => {
    post.addEventListener("click", (e) => {
      let countEl = post.querySelector(".like-count");
      const upLike = e.target.closest(".like-up");
      if (upLike) {
        e.preventDefault();
        reactToPost(+post.dataset.id, 1, countEl);
        return;
      }
      const downLike = e.target.closest(".like-down");
      if (downLike) {
        e.preventDefault();
        reactToPost(+post.dataset.id, -1, countEl);
        return;
      }
    });
  });
}

function countEffect(countEl, isSucces) {
  if (isSucces) {
    countEl.style.color = "green";
  } else {
    countEl.style.color = "red";
  }

  setTimeout(() => {
    countEl.style.color = "";
  }, 1000);
}
