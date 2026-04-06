import { Injectable, Logger } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { catchError, firstValueFrom } from 'rxjs';

@Injectable()
export class WebhookService {
  private readonly logger = new Logger(WebhookService.name);

  constructor(private readonly httpService: HttpService) {}

  async notifyGoService(eventData: any) {
    const webhookUrl = process.env.GO_WEBHOOK_URL || 'http://localhost:8080/webhook';
    try {
      this.logger.log(`Firing webhook to ${webhookUrl}`);
      const response = await firstValueFrom(
        this.httpService.post(webhookUrl, eventData).pipe(
          catchError((error) => {
            this.logger.error(`Webhook failed: ${error.message}`);
            throw 'An error happened!';
          }),
        ),
      );
      return response.data;
    } catch (e) {
      this.logger.error(`Failed to send webhook to Go service`);
    }
  }
}
