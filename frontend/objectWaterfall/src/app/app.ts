import { Component, inject, OnDestroy, OnInit, output, signal, WritableSignal } from '@angular/core';
import { WorkerSettings } from "./worker-settings/worker-settings";
import { SeedData } from "./seed-data/seed-data";
import { StartWorker } from "./start-worker/start-worker";
import { WorkerLogs } from "./worker-logs/worker-logs";
import { WorkersList } from "./workers-list/workers-list";
import { WorkerItemModel } from './models/worker/worker-item';
import { HttpClient } from '@angular/common/http';
import { WorkerRealtimeLogs } from './services/realtime/web-sockets.service';
import { Subscription } from 'rxjs';
import { environment } from './environments/environments';
import { LogModel } from './models/worker/worker-log';

class ResultMap {
  result: {id: number, name: string}[] = []
}

const DUMMY_LOGS: LogModel[] = [
    // {
    //   Log: "test1",
    //   SuccessAttemptsCount: 1,
    //   FailedAttemptsCount: 0
    // },
    // {
    //   Log: "test2",
    //   SuccessAttemptsCount: 2,
    //   FailedAttemptsCount: 0
    // },
    // {
    //   Log: "test3",
    //   SuccessAttemptsCount: 3,
    //   FailedAttemptsCount: 754
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // },
    // {
    //   Log: "test4",
    //   SuccessAttemptsCount: 0,
    //   FailedAttemptsCount: 5
    // }
  ]

@Component({
  selector: 'app-root',
  imports: [WorkerSettings, SeedData, StartWorker, WorkerLogs, WorkersList],
  templateUrl: './app.html',
  styleUrls: ['./app.css',
    '../assets/styles/settings-controls.css'
  ]
})
export class App implements OnInit, OnDestroy {
  private http = inject(HttpClient);
  private websocketService = inject(WorkerRealtimeLogs)
  private subscription!: Subscription
  private receivedMessages: LogModel[] = []

  runningWorkers = signal<WorkerItemModel[]>([])
  existingWorkers = signal<WorkerItemModel[]>([])
  isLoading = signal<boolean>(false)
  errorMessage = signal<string | null>(null)
  workerLogs = signal<LogModel[]>([])
  selectedForStartWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  selectedRunningWorker = signal<WorkerItemModel>(new WorkerItemModel(0, ""))
  isRunningWorkerSet = signal<boolean>(false);

  ngOnDestroy(): void {
    this.subscription.unsubscribe()
    this.websocketService.close()
  }

  ngOnInit(): void {
    this.getRunningWorkers()
    this.getExistingWorkers()
    this.websocketService.startConnection(environment.baseAddress + 'logsWs')
    this.subscription = this.websocketService.messages$.subscribe({
      next: (msg: LogModel) => {
        let updated: LogModel[] = []
        if (this.receivedMessages.length >= 10){
          this.receivedMessages.pop()
          this.receivedMessages.push(new LogModel(msg))
          updated = [...this.receivedMessages].reverse();
        } else {
          this.receivedMessages.push(new LogModel(msg))
          updated = [...this.receivedMessages].reverse();
        }
        this.workerLogs.set(updated)
      },
      error: err => this.errorMessage = err,
      complete: () => console.log('Socket close')
    })
  }

  sendPing(): void {
    this.websocketService.send({ type: 'PING', timestamp: Date.now() });
  }

  getRunningWorkers() {
    this.getWorkers(this.runningWorkers, 'getRunningWorkers')  
  }

  getExistingWorkers() {
    this.getWorkers(this.existingWorkers, 'getWorkers')
    console.log(this.existingWorkers())
  }

  getWorkers(workerCollection: WritableSignal<WorkerItemModel[]>, path: string){
    this.isLoading.set(true)
    this.http.get<ResultMap>(environment.baseAddress + path).subscribe({
      next: response => {
        let workers = []
        if (response.result !== null){
        for (let i = 0; i < response.result.length; i++) {
                  workers[i] = new WorkerItemModel(response.result[i].id, response.result[i].name);
                }
        workerCollection.set(workers)
        }
        this.isLoading.set(false)
      },
      error: err => {
        this.errorMessage.set(err.error.error)
        this.isLoading.set(false)
      }
    });
  }

  onSelectedRunningWorker(id: number) {
    this.websocketService.send({"workerId" : id})
    this.selectedRunningWorker.set(this.runningWorkers().find(x => x.id == id)!)
    this.isRunningWorkerSet.set(true)
  }
}
