document.addEventListener("DOMContentLoaded", function () {
  const form = document.getElementById("resetForm");
  const password = document.getElementById("password");
  const passwordConfirm = document.getElementById("password_confirmation");
  const submitButton = form.querySelector('button[type="submit"]');
  const requirementsList = document.createElement("div");
  requirementsList.id = "password-requirements";
  requirementsList.innerHTML = `
      <div class="text-muted mb-2">パスワード要件:</div>
      <ul class="text-muted pl-4" style="font-size: 0.9em;">
          <li id="req-length">8文字以上</li>
          <li id="req-no-space">スペース禁止</li>
          <li id="req-no-fullwidth">全角文字禁止</li>
          <li id="req-format">半角英数字と記号のみ使用可能</li>
          <li id="req-match">パスワードの一致</li>
      </ul>
  `;
  requirementsList.style.display = "none";
  password.parentNode.insertAdjacentElement("afterend", requirementsList);

  function validatePassword(pwd) {
    // 要件のチェック
    const isLengthValid = pwd.length >= 8;
    const hasNoSpace = !/\s/.test(pwd);
    const hasNoFullWidth = !/[^\x00-\x7F]/.test(pwd);
    const isFormatValid = /^[A-Za-z0-9!@#$%^&*()_+\-=[\]{};:"\\|,.<>/?]+$/.test(
      pwd
    );

    // 要件表示の更新
    document.getElementById("req-length").style.color = isLengthValid
      ? "green"
      : "red";
    document.getElementById("req-no-space").style.color = hasNoSpace
      ? "green"
      : "red";
    document.getElementById("req-no-fullwidth").style.color = hasNoFullWidth
      ? "green"
      : "red";
    document.getElementById("req-format").style.color = isFormatValid
      ? "green"
      : "red";

    // すべての要件を満たしているかチェック
    return isLengthValid && hasNoSpace && hasNoFullWidth && isFormatValid;
  }

  function checkSubmitButtonState() {
    const pwd = password.value;
    const pwdConfirm = passwordConfirm.value;

    const isPasswordValid = validatePassword(pwd);
    const isPasswordMatching = pwd === pwdConfirm;

    // パスワード一致の要件表示
    document.getElementById("req-match").style.color = isPasswordMatching
      ? "green"
      : "red";

    // 要件リストの表示/非表示
    requirementsList.style.display = pwd ? "block" : "none";

    // 送信ボタンの有効/無効の切り替え
    const isSubmitEnabled = isPasswordValid && isPasswordMatching;
    submitButton.disabled = !isSubmitEnabled;
  }

  // パスワード入力時の動的バリデーション
  password.addEventListener("input", function () {
    checkSubmitButtonState();
  });

  // パスワード確認入力時の検証
  passwordConfirm.addEventListener("input", function () {
    checkSubmitButtonState();
  });

  // 初期状態で送信ボタンを無効化
  submitButton.disabled = true;
});
