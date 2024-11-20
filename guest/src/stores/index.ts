import { defineStore } from "pinia";
import { parseToken } from "@/utils/parseToken";
import { UserInfoApi } from "@/api/user";
import { ElMessage } from "element-plus";

interface userInfoType {
  exp: number;
  Uid: number;
  token: string;
  avatar: string;
}

const userInfo: userInfoType = {
  exp: 0,
  Uid: 0,
  token: "",
  avatar: "",
};

export const useStore = defineStore("counter", {
  state: () => {
    return {
      userInfo: userInfo,
    };
  },

  actions: {
    async setToken(token: string) {
      try {
        const payload = parseToken(token);
        this.userInfo.token = token;
        this.userInfo.exp = payload.exp;
        this.userInfo.Uid = payload.Uid;
        // 去拿用户的信息
        let res = await UserInfoApi();
        if (res.code) {
          ElMessage.error(res.msg);
          return;
        }
        this.userInfo.avatar = res.data.avatar;
        // 持久化存储
        this.saveToken();
      } catch (e) {
        console.log(e);
      }
    },
    saveToken() {
      try {
        localStorage.setItem("UserInfo", JSON.stringify(this.userInfo));
      } catch (error) {
        console.error(error);
      }
    },
    loadToken() {
      const val = localStorage.getItem("UserInfo");
      if (!val) {
        // 没有登录，或者登录失效
        return;
      }
      try {
        this.userInfo = JSON.parse(val);
      } catch (e) {
        localStorage.removeItem("UserInfo");
        return;
      }
    },
  },

  getters: {
    isLogin(): boolean {
      try {
        return this.userInfo.token != "";
      } catch (e) {
        console.log(e);
      }
    },
  },
});
