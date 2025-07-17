// frontend/js/api.js

const API_BASE_URL = "https://your-api-gateway-url.execute-api.ap-northeast-1.amazonaws.com/prod"; // 後で差し替え

/**
 * GETリクエストを送信
 * @param {string} path APIパス（例: '/careplans'）
 * @returns {Promise<any>} JSONレスポンス
 */
export async function getRequest(path) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      // 認証トークンがあればここに追加
    },
  });
  if (!response.ok) {
    throw new Error(`GET ${path} failed: ${response.status}`);
  }
  return response.json();
}

/**
 * POSTリクエストを送信
 * @param {string} path APIパス
 * @param {object} data 送信データ
 * @returns {Promise<any>} JSONレスポンス
 */
export async function postRequest(path, data) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // 認証トークンがあればここに追加
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`POST ${path} failed: ${response.status}`);
  }
  return response.json();
}

// ここにPUT, DELETEも同様に実装可能
