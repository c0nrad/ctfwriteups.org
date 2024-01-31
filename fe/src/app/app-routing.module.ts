import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { SearchComponent } from "./pages/search/search.component";
import { WriteupViewComponent } from "./pages/writeup-view/writeup-view.component";
import { WriteupEditComponent } from "./pages/writeup-edit/writeup-edit.component";
import { LoginComponent } from "./pages/login/login.component";
import { CtfListComponent } from "./pages/ctf-list/ctf-list.component";
import { CtfEditComponent } from "./pages/ctf-edit/ctf-edit.component";
import { NewestComponent } from "./pages/newest/newest.component";
import { TopComponent } from "./pages/top/top.component";
import { AboutComponent } from "./pages/about/about.component";
import { NewsletterComponent } from "./pages/newsletter/newsletter.component";
import { UnsubscribeComponent } from "./pages/unsubscribe/unsubscribe.component";
import { CtfViewComponent } from "./pages/ctf-view/ctf-view.component";

const routes: Routes = [
  { path: "", component: TopComponent },
  { path: "newest", component: NewestComponent },
  { path: "search", component: SearchComponent },
  { path: "writeups/:id", component: WriteupViewComponent },
  { path: "writeups/:id/edit", component: WriteupEditComponent },
  { path: "submit", component: WriteupEditComponent },

  { path: "ctfs", component: CtfListComponent },
  { path: "ctfs/:id/edit", component: CtfEditComponent },
  { path: "ctfs/:id", component: CtfViewComponent },

  { path: "newsletter", component: NewsletterComponent },
  { path: "unsubscribe", component: UnsubscribeComponent },

  { path: "about", component: AboutComponent },

  { path: "login", component: LoginComponent },

  { path: "**", redirectTo: "top" },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
