import { Controller, Get, Post, Param, Delete, Body, UseGuards, Request } from '@nestjs/common';
import { EventsService } from './events.service';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { RolesGuard } from '../auth/roles.guard';
import { Roles } from '../auth/roles.decorator';
import { Role } from '@prisma/client';

@Controller('api/v1/events')
@UseGuards(JwtAuthGuard, RolesGuard)
export class EventsController {
  constructor(private readonly eventsService: EventsService) {}

  @Post()
  @Roles(Role.ADMIN) // Only admins can create events
  create(@Body() body: any) {
    return this.eventsService.create({
      title: body.title,
      description: body.description,
      date: new Date(body.date),
    });
  }

  @Get()
  findAll() {
    return this.eventsService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.eventsService.findOne(id);
  }

  @Post(':id/join')
  joinEvent(@Param('id') id: string, @Request() req) {
    return this.eventsService.joinEvent(req.user.id, id);
  }

  @Delete(':id')
  @Roles(Role.ADMIN) // Bonus A details: admin roles manage events
  remove(@Param('id') id: string) {
    return this.eventsService.remove(id);
  }
}
