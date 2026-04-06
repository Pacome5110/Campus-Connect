import { Resolver, Query, Args, Mutation, ResolveField, Parent } from '@nestjs/graphql';
import { EventsService } from './events.service';
import { Event, User } from '../models.graphql';
import { UseGuards, Request } from '@nestjs/common';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { PrismaService } from '../prisma/prisma.service';

@Resolver(() => Event)
export class EventsResolver {
  constructor(
    private eventsService: EventsService,
    private prisma: PrismaService,
  ) {}

  @Query(() => [Event])
  @UseGuards(JwtAuthGuard)
  async getEvents() {
    return this.eventsService.findAll();
  }

  @Query(() => Event)
  @UseGuards(JwtAuthGuard)
  async getEvent(@Args('id') id: string) {
    return this.eventsService.findOne(id);
  }

  @ResolveField(() => [User])
  async participants(@Parent() event: Event) {
    const records = await this.prisma.eventParticipant.findMany({
      where: { eventId: event.id },
      include: { user: true },
    });
    return records.map(r => r.user);
  }
}
