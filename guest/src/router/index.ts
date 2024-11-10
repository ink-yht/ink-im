import { createRouter, createWebHistory } from "vue-router";
import { useStore } from "@/stores";
import { ElMessage } from "element-plus";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // login
    {
      path: "/login",
      name: "login",
      component: () => import("@/views/login/login.vue"),
    },

    // signup
    {
      path: "/signup",
      name: "signup",
      component: () => import("@/views/login/signup.vue"),
    },

    // web
    {
      path: "/",
      name: "web",
      component: () => import("@/views/web/index.vue"),
      children: [
        // contacts
        {
          path: "",
          name: "contacts", // 联系人
          component: () => import("@/views/web/contacts/index.vue"),
          children: [
            {
              path: "",
              name: "user_list",
              component: () => import("@/views/web/contacts/user_list.vue"),
            },
            {
              path: "welcome",
              name: "welcome",
              component: () => import("@/views/web/contacts/welcome.vue"),
            },
          ],
        },

        // session
        {
          path: "session",
          name: "session",
          component: () => import("@/views/web/session/index.vue"),
        },

        // info
        {
          path: "info",
          name: "info",
          component: () => import("@/views/web/info/index.vue"),
        },

        // notice
        {
          path: "notice",
          name: "notice",
          component: () => import("@/views/web/notice/index.vue"),
        },
      ],
      meta: {
        isLogin: true, // 需要登陆验证
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  console.log("to:", to);
  console.log("from:", from);

  if (to.meta.isLogin === true) {
    // 查询有没有登录
    const store = useStore();
    if (!store.isLogin) {
      // 没有登录，跳转到登录页面
      const redirectPath = from.path !== "/" ? from.path : to.path;
      router.push({
        name: "login",
        query: {
          redirect_url: redirectPath,
        },
      });
      console.log("from.path", from.path);
      ElMessage.warning("请登录");
      return;
    }
  }
  next();
});

export default router;
