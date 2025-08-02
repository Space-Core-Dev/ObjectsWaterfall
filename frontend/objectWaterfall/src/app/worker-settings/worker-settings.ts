import { Component } from '@angular/core';
import { signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { inject } from '@angular/core';

class Settings {
  workerName = ""
  timer = 30
  requestDelay = 1
  random = false
  writesNumberToSend = 10
  totalToSend = 1000
  StopWhenTableEnds = false
}

@Component({
  selector: 'app-worker-settings',
  imports: [FormsModule],
  templateUrl: './worker-settings.html',
  styleUrl: './worker-settings.css'
})
export class WorkerSettings {
  newSettings = signal<Settings>(new Settings())
  errorMessage = signal<string | null>(null)
  isLoading = signal<boolean>(false)
  private http = inject(HttpClient);

  onAdd(){
    this.errorMessage.set(null)
    this.isLoading.set(true)
    this.sendSettings()
  }

  updateRandom(event: Event) {
    const isChecked = (event.target as HTMLInputElement).checked;
    this.newSettings.update(settings => ({
      ...settings,
      random: isChecked
    }));
  }

  updateEndData(event: Event) {
    const isChecked = (event.target as HTMLInputElement).checked;
    this.newSettings.update(settings => ({
      ...settings,
      StopWhenTableEnds: isChecked
    }));
  }

  sendSettings() {
    const payload = this.newSettings();
    console.log(payload)
    this.http.post('http://localhost:8888/add', payload).subscribe({
      next: response => {
        this.isLoading.set(false)
      },
      error: err => {
        this.errorMessage.set(err)
        this.isLoading.set(false)
      }
    });
  }
}
