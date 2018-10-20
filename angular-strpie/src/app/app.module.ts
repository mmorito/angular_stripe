import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { ComponentsComponent } from './components/components.component';
import { StripeComponent } from './components/stripe/stripe.component';

@NgModule({
  declarations: [
    AppComponent,
    ComponentsComponent,
    StripeComponent
  ],
  imports: [
    BrowserModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
