import { Component, inject, output, signal } from '@angular/core';
import { WorkerSettings } from "./worker-settings/worker-settings";
import { SeedData } from "./seed-data/seed-data";
import { StartWorker } from "./start-worker/start-worker";
import { WorkerLogs } from "./worker-logs/worker-logs";
import { WorkersList } from "./workers-list/workers-list";
import { WorkerItemModel } from './models/worker/worker-item';
import { HttpClient } from '@angular/common/http';

class ResultMap {
  result: string[] = []
}

@Component({
  selector: 'app-root',
  imports: [WorkerSettings, SeedData, StartWorker, WorkerLogs, WorkersList],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  private http = inject(HttpClient);
  workers = signal<WorkerItemModel[]>([])
  isLoading = signal<boolean>(false)
  errorMessage = signal<string | null>(null)

  constructor(){
    this.getWorkers()
  }

  getWorkers() {
    this.isLoading.set(true)
    this.http.get<ResultMap>('http://localhost:8888/getWorkers').subscribe({
      next: response => {
        let workerNames = []
        for (let i = 0; i < response.result.length; i++) {
          workerNames[i] = new WorkerItemModel(response.result[i]);
        }
        this.workers.set(workerNames)
        this.isLoading.set(false)
      },
      error: err => {
        this.errorMessage.set(err.error.error)
        this.isLoading.set(false)
      }
    });
  }
}
