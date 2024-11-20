<script setup lang="ts">
import { onMounted, reactive, watch } from "vue";
import CryptoJS from "crypto-js";
import { Hide, View, Message, Lock } from "@element-plus/icons-vue";
import Qq_color from "@/components/im_login/qq_color.vue";
import { ElMessage } from "element-plus";
import { LoginApi, type LoginRequest } from "@/api/user";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();

// 表单校验
const rules = reactive({
  email: [{ required: true, message: "Please input email", trigger: "blur" }],
  password: [
    { required: true, message: "Please input password", trigger: "blur" },
    { min: 8, max: 16, message: "Length should be 8 to 15", trigger: "blur" },
  ],
});

// 定义表单字段
const form = reactive<LoginRequest>({
  email: "",
  password: "",
  rememberMe: false,
  showPassword: false,
});

onMounted(() => {
  // 检查 localStorage 中是否有凭证
  if (form.email) {
    checkRememberMe();
  }
});
watch(
  () => form.email,
  () => {
    checkRememberMe();
  },
);

// 小眼睛查看密码
const togglePasswordVisibility = () => {
  form.showPassword = !form.showPassword;
};

// 记住密码
const checkRememberMe = () => {
  if (form.email) {
    const rememberMe =
      localStorage.getItem(`rememberMe_${form.email}`) === "true";
    if (rememberMe) {
      const encryptedEmail = localStorage.getItem(`email_${form.email}`);
      const encryptedPassword = localStorage.getItem(`password_${form.email}`);

      if (encryptedEmail && encryptedPassword) {
        // 解密凭证
        const bytesEmail = CryptoJS.AES.decrypt(encryptedEmail, "secret-key");
        const bytesPassword = CryptoJS.AES.decrypt(
          encryptedPassword,
          "secret-key",
        );
        const email = bytesEmail.toString(CryptoJS.enc.Utf8);
        const password = bytesPassword.toString(CryptoJS.enc.Utf8);

        // 自动填充表单
        form.email = email;
        form.password = password;
        form.rememberMe = true;
      }
    }
  }
};

// 表单登录
const handleLogin = async () => {
  if (form.email && form.password) {
    // 实现记住密码
    if (form.rememberMe) {
      // 加密并存储凭证
      const encryptedEmail = CryptoJS.AES.encrypt(
        form.email,
        "secret-key",
      ).toString();
      const encryptedPassword = CryptoJS.AES.encrypt(
        form.password,
        "secret-key",
      ).toString();
      localStorage.setItem(`email_${form.email}`, encryptedEmail);
      localStorage.setItem(`password_${form.email}`, encryptedPassword);
      localStorage.setItem(`rememberMe_${form.email}`, "true");
    } else {
      // 清除凭证
      localStorage.removeItem(`email_${form.email}`);
      localStorage.removeItem(`password_${form.email}`);
      localStorage.removeItem(`rememberMe_${form.email}`);
    }

    // 向后端发起请求
    let res = await LoginApi(form);
    // code 不为 0
    if (res.code) {
      ElMessage.error(res.msg);
      return;
    }
    ElMessage.success(res.msg);

    // 先拿 redirect_url,如果有，就跳转到这里
    const redirectUrl = router.currentRoute.value.query.redirect_url;
    if (redirectUrl) {
      await router.push({
        path: redirectUrl as string,
      });
      return;
    }
    await router.push({
      name: "web",
    });
  } else {
    ElMessage.error("please input it to completion");
  }
};
</script>

<template>
  <div class="im_login">
    <div class="banner">
      <qq_color></qq_color>
    </div>
    <div class="login_form">
      <el-form :model="form" :rules="rules">
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="邮箱"
            :prefix-icon="Message"
          >
          </el-input>
        </el-form-item>
        <el-form-item class="item_password" prop="password">
          <!-- 使用 :type 动态绑定 input 类型 -->
          <el-input
            :type="form.showPassword ? 'text' : 'password'"
            v-model="form.password"
            placeholder="密码"
            :prefix-icon="Lock"
          >
            <!-- 添加一个 slot 来放置密码可见性切换按钮 -->
            <template #suffix>
              <i
                v-if="form.showPassword"
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <Hide />
                </el-icon>
              </i>
              <i
                v-else
                style="cursor: pointer; font-size: 18px"
                @click="togglePasswordVisibility"
              >
                <el-icon>
                  <View />
                </el-icon>
              </i>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item class="item_action">
          <el-checkbox v-model="form.rememberMe">记住密码</el-checkbox>
          <router-link
            :underline="false"
            type="primary"
            class="right"
            to="/signup"
            >注册账号</router-link
          >
        </el-form-item>
        <el-form-item>
          <el-button style="width: 100%" type="primary" @click="handleLogin"
            >登录</el-button
          >
        </el-form-item>
      </el-form>
      <div class="other_login">
        <div class="label">第三方登录</div>
        <div class="icons">
          <i class="iconfont icon-QQ"></i>
          <i class="iconfont icon-weixin"></i>
          <i class="iconfont icon-shoujitongxun"></i>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.im_login {
  width: 500px;
  height: 406px;
  background-color: white;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 0 5px 3px rgba(0, 0, 0, 0.05);

  .banner {
    height: 140px;
    width: 100%;
    background-color: aquamarine;
  }

  .login_form {
    padding: 20px 60px;
  }

  .item_password {
    margin-bottom: 10px;
  }

  .item_action {
    margin-bottom: 4px;
    position: relative;

    .right {
      position: absolute;
      right: 0.1px;
      color: #606266;
      cursor: pointer;
      text-decoration: none;

      &:hover {
        color: #409eff;
      }
    }
  }

  .other_login {
    display: flex;
    flex-direction: column;
    align-items: center;

    .label {
      font-size: 14px;
      color: #aaa8a8;
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;

      &::before,
      &::after {
        width: 35%;
        height: 1px;
        background-color: #aaa8a8;
        content: "";
        display: inline-flex;
      }
    }

    .icons {
      display: flex;
      justify-content: center;
      align-items: center;
      margin-top: 6px;

      i {
        font-size: 36px;
        cursor: pointer;
        margin: 0 20px;
      }
    }
  }
}
</style>
