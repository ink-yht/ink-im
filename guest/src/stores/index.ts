import { defineStore } from "pinia";
import { parseToken } from "@/utils/parseToken";

interface UserInfoType {
  exp: number;
  Uid: number;
  token: string;
  avatar: string;
}

const initialUserInfo: UserInfoType = {
  exp: 0,
  Uid: 0,
  token: "",
  avatar: "",
};

export const useStore = defineStore("mainStore", {
  state: () => ({
    UserInfo: { ...initialUserInfo }, // 使用扩展运算符复制初始状态
  }),

  actions: {
    async setToken(token: string) {
      try {
        const payload = parseToken(token) as Partial<UserInfoType>; // 确保 parseToken 返回的是 Partial<UserInfoType>

        if (payload.exp !== undefined) this.UserInfo.exp = payload.exp;
        if (payload.Uid !== undefined) this.UserInfo.Uid = payload.Uid;
        this.UserInfo.token = token;

        // 去拿用户的基本信息
        let res = await this.saveToken();
      } catch (error) {
        console.error("Failed to set token:", error);
      }
    },

    saveToken() {
      try {
        localStorage.setItem("UserInfo", JSON.stringify(this.UserInfo));
      } catch (error) {
        console.error("Failed to save token:", error);
      }
    },

    loadToken() {
      try {
        const val = localStorage.getItem("UserInfo");
        if (!val) {
          // 没有登录，或者登录失效
          return;
        }
        this.userInfo = JSON.parse(val);
      } catch (error) {
        localStorage.removeItem("UserInfo");
        console.error("Failed to load token:", error);
      }
    },
  },

  getters: {},
});
