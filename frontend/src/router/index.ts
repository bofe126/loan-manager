import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import Home from '../views/Home.vue';
import LoanListPage from '../views/LoanListPage.vue';
import AddLoan from '../views/AddLoan.vue';
import EditLoan from '../views/EditLoan.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/loans',
    name: 'LoanList',
    component: LoanListPage,
  },
  {
    path: '/add',
    name: 'AddLoan',
    component: AddLoan,
  },
  {
    path: '/edit/:id',
    name: 'EditLoan',
    component: EditLoan,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
