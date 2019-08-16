import { Component } from '@angular/core';

@Component({
  selector: 'account-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.less'],
})

export class AccountLayoutComponent {
  links = [
    {
      title: '帮助',
      href: '',
    },
    {
      title: '隐私',
      href: '',
    },
    {
      title: '条款',
      href: '',
    },
  ];

  title = 'Voice Comic';
}
