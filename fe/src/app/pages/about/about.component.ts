import { Component } from "@angular/core";

@Component({
  selector: "app-about",
  template: `
    <div class="container" style="padding-top: 65px">
      <h3>CTFWriteups.org</h3>

      <p>CTFWriteups.org is a place to share and learn from CTF Writeups</p>

      <h3>F.A.Q.</h3>

      <h4>Why did you build this?</h4>
      <p>
        The idea started <a href="https://twitter.com/c0nrad_jr/status/1726216778077053363">here</a> and grew. CTF writeups are the best way to learn and I wanted a tool to help organize and share
        them.
      </p>

      <h4>How can I report a bug or feature?</h4>

      <p>Join the <a href="https://discord.gg/srDFwyHmVS">discord.</a></p>

      <h4>Features</h4>
      <ul>
        <li>Quickly view all writeups related to a CTF</li>
        <li>View the top rated writeups</li>
        <li>Search previous writeups for ideas during CTFs</li>
        <li>Get feedback on your writeups through comments/upvotes</li>
        <li>Receive a weekly digest with the latest and greatest writeups</li>
      </ul>

      <h4>Who made this?</h4>
      <p><a href="https://twitter.com/c0nrad_jr">c0nrad</a> of the <a href="https://youtube.com/@SloppyJoePirates">Sloppy Joe Pirates</a></p>
    </div>
  `,
  styles: [],
})
export class AboutComponent {}
