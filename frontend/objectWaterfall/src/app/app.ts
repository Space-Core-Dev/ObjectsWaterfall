import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { WorkerSettings } from "./worker-settings/worker-settings";

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, WorkerSettings],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected title = 'objectWaterfall';
}
