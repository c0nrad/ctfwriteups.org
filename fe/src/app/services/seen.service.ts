import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ConfigService } from "./config.service";

export interface Seen {
  id: string;
  ts: number;

  writeupID: string;
  userID: string;
}

@Injectable({
  providedIn: "root",
})
export class SeenService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getSeens(): Observable<Seen[]> {
    return this.http.get<Seen[]>(`${this.config.getOrigin()}/api/v1/users/me/seens`);
  }

  newSeen(seen: Seen): Observable<Seen> {
    return this.http.post<Seen>(`${this.config.getOrigin()}/api/v1/seens`, seen);
  }

  deleteSeen(seen: Seen): Observable<void> {
    return this.http.delete<void>(`${this.config.getOrigin()}/api/v1/seens/${seen.id}`);
  }
}
