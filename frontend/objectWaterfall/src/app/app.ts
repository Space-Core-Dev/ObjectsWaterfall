import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { WorkerSettings } from "./worker-settings/worker-settings";
import { SeedData } from "./seed-data/seed-data";
import { StartWorker } from "./start-worker/start-worker";
import { WorkerLogs } from "./worker-logs/worker-logs";

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, WorkerSettings, SeedData, StartWorker, WorkerLogs],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected title = 'objectWaterfall';
}
