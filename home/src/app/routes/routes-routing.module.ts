import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {SimpleGuard} from '@delon/auth';
import {environment} from '@env/environment';
// layout
import {LayoutDefaultComponent} from '../layout/default/default.component';
import {LayoutFullScreenComponent} from '../layout/fullscreen/fullscreen.component';

// dashboard pages
import {DashboardComponent} from './dashboard/dashboard.component';
// passport pages
import {UserRegisterComponent} from './passport/register/register.component';
import {UserRegisterResultComponent} from './passport/register-result/register-result.component';
// single pages
import {CallbackComponent} from './callback/callback.component';
import {UserLockComponent} from './passport/lock/lock.component';

// account
import {AccountLoginComponent} from '../views/account/login/login.component'
import {AccountLayoutComponent} from "../views/account/layout/layout.component";

const routes: Routes = [
  {
    path: '',
    component: LayoutDefaultComponent,
    canActivate: [SimpleGuard],
    children: [
      {path: '', redirectTo: 'dashboard', pathMatch: 'full'},
      {path: 'dashboard', component: DashboardComponent, data: {title: '仪表盘', titleI18n: 'dashboard'}},
      {path: 'exception', loadChildren: () => import('./exception/exception.module').then(m => m.ExceptionModule)},
      // 业务子模块
      // { path: 'widgets', loadChildren: () => import('./widgets/widgets.module').then(m => m.WidgetsModule) },
    ]
  },
  // 全屏布局
  // {
  //     path: 'fullscreen',
  //     component: LayoutFullScreenComponent,
  //     children: [
  //     ]
  // },
  // passport
  {
    path: 'passport',
    component: AccountLayoutComponent,
    children: [
      {
        path: 'login',
        component: AccountLoginComponent,
        data: {title: '登录', titleI18n: 'Login'}
      },
      {
        path: 'register',
        component: UserRegisterComponent,
        data: {title: '注册', titleI18n: 'Register'}
      },
      {
        path: 'register-result',
        component: UserRegisterResultComponent,
        data: {title: '注册结果', titleI18n: 'pro-register-result'}
      },
      {
        path: 'lock',
        component: UserLockComponent,
        data: {title: '锁屏', titleI18n: 'lock'}
      },
    ]
  },
  // 单页不包裹Layout
  {path: 'callback/:type', component: CallbackComponent},

  // Account
  {
    path: 'account',
    component: AccountLayoutComponent,
    children: [
      {
        path: 'login',
        component: AccountLoginComponent,
        data: {title: '登录', titleI18n: 'Login'}
      }
      ]
  },

  {path: '**', redirectTo: 'exception/404'},
];

@NgModule({
  imports: [
    RouterModule.forRoot(
      routes, {
        useHash: environment.useHash,
        // NOTICE: If you use `reuse-tab` component and turn on keepingScroll you can set to `disabled`
        // Pls refer to https://ng-alain.com/components/reuse-tab
        scrollPositionRestoration: 'top',
      }
    )],
  exports: [RouterModule],
})
export class RouteRoutingModule {
}
