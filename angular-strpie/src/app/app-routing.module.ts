import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { StripeComponent } from './components/stripe/stripe.component';

const routes: Routes = [
  // { path: '', component: ChatComponent },
  { path: 'stripe', component: StripeComponent }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}
