document.addEventListener("DOMContentLoaded", function () {
  const passwordForm = document.getElementById("password-form");
  const currentPassword = document.getElementById("current-password");
  const newPassword = document.getElementById("new-password");
  const confirmPassword = document.getElementById("confirm-password");
  const submitButton = passwordForm.querySelector('button[type="submit"]');
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
          <li id="req-different">現在のパスワードと異なる</li>
      </ul>
  `;
  requirementsList.style.display = "none";
  newPassword.parentNode.insertAdjacentElement("afterend", requirementsList);

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
    const current = currentPassword.value;
    const pwd = newPassword.value;
    const pwdConfirm = confirmPassword.value;

    const isPasswordValid = validatePassword(pwd);
    const isPasswordMatching = pwd === pwdConfirm;
    const isDifferentFromCurrent = pwd !== current;

    // 要件表示の更新
    document.getElementById("req-match").style.color = isPasswordMatching
      ? "green"
      : "red";
    document.getElementById("req-different").style.color =
      isDifferentFromCurrent ? "green" : "red";

    // 要件リストの表示/非表示
    requirementsList.style.display = pwd ? "block" : "none";

    // 送信ボタンの有効/無効の切り替え
    const isSubmitEnabled =
      current &&
      isPasswordValid &&
      isPasswordMatching &&
      isDifferentFromCurrent;
    submitButton.disabled = !isSubmitEnabled;
  }

  // 各入力フィールドのイベントリスナー
  currentPassword.addEventListener("input", checkSubmitButtonState);
  newPassword.addEventListener("input", checkSubmitButtonState);
  confirmPassword.addEventListener("input", checkSubmitButtonState);

  // プロフィール設定フォームの送信
  document
    .getElementById("profile-form")
    .addEventListener("submit", function (e) {
      e.preventDefault();
      // TODO: プロフィール更新の処理
      alert("プロフィールを更新しました");
    });

  // 通知設定フォームの送信
  document
    .getElementById("notification-form")
    .addEventListener("submit", function (e) {
      e.preventDefault();
      // TODO: 通知設定の更新処理
      alert("通知設定を更新しました");
    });

  // アカウント削除ボタンのイベント
  document
    .getElementById("delete-account")
    .addEventListener("click", function () {
      if (
        confirm("本当にアカウントを削除しますか？この操作は取り消せません。")
      ) {
        // TODO: アカウント削除の処理
        window.location.href = "/settings/account/delete";
      }
    });

  // 初期状態で送信ボタンを無効化
  submitButton.disabled = true;
});
