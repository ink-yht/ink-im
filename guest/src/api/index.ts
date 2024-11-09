import axios from "axios";
import { ElMessage } from "element-plus";
import { parseToken } from "@/utils/parseToken";
import { useStore } from "@/stores";

const store = useStore();

export const useAxios = axios.create({
  baseURL: "http://127.0.0.1:8081",
});

export interface baseResponse<T> {
  code: number;
  data: T;
  msg: string;
}
export interface listDataType<T> {
  count: number;
  list: T[];
}

useAxios.interceptors.request.use((config) => {
  // const store = useStore();
  // config.headers["token"] = store.userInfo.token;
  return config;
});
useAxios.interceptors.response.use(
  (response) => {
    if (response.status !== 200) {
      // 失败的
      console.log("服务失败", response.status);
      ElMessage.error(response.statusText);
      return Promise.reject(response.statusText);
    }

    // 获取 token
    const newToken = response.headers["x-jwt-token"];
    store.setToken(newToken);
    return response.data;
  },
  (err) => {
    console.log("服务错误", err);
    ElMessage.error(err.message);
    return Promise.reject(err.message);
  },
);
