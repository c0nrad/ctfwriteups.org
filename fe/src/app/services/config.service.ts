import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { environment } from "src/environments/environment";

@Injectable({
  providedIn: "root",
})
export class ConfigService {
  constructor(private http: HttpClient) {}

  getAssetsOrigin(): string {
    if (this.getHost().indexOf("assets") > -1) {
      return this.getOrigin();
    } else if (this.getHost().indexOf("ctfwriteups") > -1) {
      return "https://assets." + this.getHost();
    } else {
      return this.getOrigin();
    }
  }

  getHost(): string {
    return window.location.host;
  }

  getOrigin(): string {
    return window.location.origin;
  }

  getWSOrigin(): string {
    if (environment.production) {
      return this.getHost();
    } else {
      return "localhost:8080";
    }
  }
}
