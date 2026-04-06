import { Module } from '@nestjs/common';
import { EventsService } from './events.service';
import { EventsController } from './events.controller';
import { EventsResolver } from './events.resolver';

@Module({
  controllers: [EventsController],
  providers: [EventsService, EventsResolver],
  exports: [EventsService]
})
export class EventsModule {}
