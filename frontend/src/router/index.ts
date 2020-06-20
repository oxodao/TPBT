import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'

import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import CallbackTwitch from '../views/CallbackTwitch.vue'
import Game from '../views/Game.vue'

Vue.use(VueRouter)

  const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/auth/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/auth/callback',
    name: 'Logged-in',
    component: CallbackTwitch
  },
    {
      path: '/game',
      name: 'Game',
      component: Game
    }
  /*{
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" * / '../views/About.vue')
  }*/
]

const router = new VueRouter({
  routes
})

export default router
