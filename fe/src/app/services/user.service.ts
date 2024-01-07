import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ConfigService } from "./config.service";

export interface User {
  id: string;
  ts: number;

  username: string;
  email: string;

  isAdmin: boolean;
  isModerator: boolean;
}

@Injectable({
  providedIn: "root",
})
export class UserService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getMe(): Observable<User> {
    return this.http.get<User>(`${this.config.getOrigin()}/api/v1/users/me`);
  }

  logout(): Observable<void> {
    return this.http.get<void>(`${this.config.getOrigin()}/logout`, {});
  }
}
