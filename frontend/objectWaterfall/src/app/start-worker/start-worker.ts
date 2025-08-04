import { Component, signal, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-start-worker',
  imports: [],
  templateUrl: './start-worker.html',
  styleUrl: './start-worker.css'
})
export class StartWorker {
   errorMessage = signal<string | null>(null)
  isMinimized = signal(false)
  private http = inject(HttpClient);

  onStart(){

  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
