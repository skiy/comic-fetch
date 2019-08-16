import { NgModule } from '@angular/core';
import { SharedModule } from '@shared';

import { LayoutDefaultComponent } from './layout/default/default.component';
import { LayoutFullScreenComponent } from './layout/fullscreen/fullscreen.component';
import { HeaderComponent } from './layout/default/header/header.component';
import { SidebarComponent } from './layout/default/sidebar/sidebar.component';
import { HeaderSearchComponent } from './layout/default/header/components/search.component';
import { HeaderNotifyComponent } from './layout/default/header/components/notify.component';
import { HeaderTaskComponent } from './layout/default/header/components/task.component';
import { HeaderIconComponent } from './layout/default/header/components/icon.component';
import { HeaderFullScreenComponent } from './layout/default/header/components/fullscreen.component';
import { HeaderStorageComponent } from './layout/default/header/components/storage.component';
import { HeaderUserComponent } from './layout/default/header/components/user.component';

import { SettingDrawerComponent } from './layout/default/setting-drawer/setting-drawer.component';
import { SettingDrawerItemComponent } from './layout/default/setting-drawer/setting-drawer-item.component';

const SETTINGDRAWER = [SettingDrawerComponent, SettingDrawerItemComponent];

const COMPONENTS = [
  LayoutDefaultComponent,
  LayoutFullScreenComponent,
  HeaderComponent,
  SidebarComponent,
  ...SETTINGDRAWER
];

const HEADERCOMPONENTS = [
  HeaderSearchComponent,
  HeaderNotifyComponent,
  HeaderTaskComponent,
  HeaderIconComponent,
  HeaderFullScreenComponent,
  HeaderStorageComponent,
  HeaderUserComponent
];

// Account layout
import { AccountLayoutComponent } from './views/account/layout/layout.component'

const ACCOUNT = [
  AccountLayoutComponent
];

@NgModule({
  imports: [SharedModule],
  entryComponents: SETTINGDRAWER,
  declarations: [
    ...COMPONENTS,
    ...HEADERCOMPONENTS,
    ...ACCOUNT
  ],
  exports: [
    ...COMPONENTS,
    ...ACCOUNT
  ]
})

export class LayoutModule { }
