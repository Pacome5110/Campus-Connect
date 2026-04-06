import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { WebhookService } from '../webhook/webhook.service';

@Injectable()
export class EventsService {
  constructor(
    private prisma: PrismaService,
    private webhookService: WebhookService,
  ) {}

  async create(data: any) {
    const event = await this.prisma.event.create({ data });
    
    // Notify Go service asynchronously
    this.webhookService.notifyGoService({ type: 'EVENT_CREATED', payload: event });
    
    return event;
  }

  async findAll() {
    return this.prisma.event.findMany({ include: { participants: true } });
  }

  async findOne(id: string) {
    const event = await this.prisma.event.findUnique({ where: { id }, include: { participants: true } });
    if (!event) throw new NotFoundException('Event not found');
    return event;
  }

  async remove(id: string) {
    return this.prisma.event.delete({ where: { id } });
  }

  async joinEvent(userId: string, eventId: string) {
    const participation = await this.prisma.eventParticipant.create({
      data: { userId, eventId },
    });
    
    // Notify Go service
    this.webhookService.notifyGoService({ type: 'USER_JOINED', payload: { userId, eventId } });
    
    return participation;
  }
}
