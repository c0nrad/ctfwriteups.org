import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";
import { HttpClientModule } from "@angular/common/http";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { SearchComponent } from "./pages/search/search.component";
import { LoginComponent } from "./pages/login/login.component";
import { WriteupViewComponent } from "./pages/writeup-view/writeup-view.component";
import { WriteupEditComponent } from "./pages/writeup-edit/writeup-edit.component";
import { FormsModule } from "@angular/forms";

import { LoadingBarHttpClientModule } from "@ngx-loading-bar/http-client";
import { LoadingBarRouterModule } from "@ngx-loading-bar/router";
import { LoadingBarModule } from "@ngx-loading-bar/core";
import { HeaderComponent } from "./components/header/header.component";
import { CtfListComponent } from "./pages/ctf-list/ctf-list.component";
import { CtfEditComponent } from "./pages/ctf-edit/ctf-edit.component";
import { NewestComponent } from "./pages/newest/newest.component";
import { TopComponent } from "./pages/top/top.component";
import { AboutComponent } from './pages/about/about.component';
import { DateAgoPipe } from './date-ago.pipe';
import { NewsletterComponent } from './pages/newsletter/newsletter.component';
import { UnsubscribeComponent } from './pages/unsubscribe/unsubscribe.component';

@NgModule({
  declarations: [AppComponent, SearchComponent, LoginComponent, WriteupViewComponent, WriteupEditComponent, HeaderComponent, CtfListComponent, CtfEditComponent, NewestComponent, TopComponent, AboutComponent, DateAgoPipe, NewsletterComponent, UnsubscribeComponent],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule, FormsModule, LoadingBarHttpClientModule, LoadingBarRouterModule, LoadingBarModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
