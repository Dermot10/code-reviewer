from sqlalchemy import Column, Integer, String, Text, ForeignKey, DateTime, Boolean, UniqueConstraint
from sqlalchemy.orm import relationship
from datetime import datetime
from .base import Base

# ------------------ USERS ------------------
class Users(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String(50), unique=True, index=True, nullable=False)
    email = Column(String(100), unique=True, index=True, nullable=False)
    fullname = Column(String(100), nullable=True)
    hashed_password = Column(String(256), nullable=False)
    is_active = Column(Boolean, default=True)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Relationships
    organisations = relationship("Organisations", back_populates="owner")
    projects = relationship("Projects", back_populates="owner")

# ------------------ ORGANISATIONS ------------------
class Organisations(Base):
    __tablename__ = "organisations"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(100), unique=True, nullable=False)
    owner_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Relationships
    owner = relationship("Users", back_populates="organisations")
    projects = relationship("Projects", back_populates="organisation")

# ------------------ PROJECTS ------------------
class Projects(Base):
    __tablename__ = "projects"

    __table_args__ = (
    UniqueConstraint("organisation_id", "name", name="uq_org_project_name"),
    )

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(100), nullable=False)
    description = Column(Text, nullable=True)
    owner_id = Column(Integer, ForeignKey("users.id"), nullable=False, index=True)
    organisation_id = Column(Integer, ForeignKey("organisations.id"), nullable=True, index=True)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Relationships
    owner = relationship("Users", back_populates="projects")
    organisation = relationship("Organisations", back_populates="projects")
    reviews = relationship("Reviews", back_populates="project", cascade="all, delete-orphan")

# ------------------ REVIEWS ------------------
class Reviews(Base):
    __tablename__ = "reviews"

    id = Column(Integer, primary_key=True, index=True)
    project_id = Column(Integer, ForeignKey("projects.id"), nullable=False, index=True)
    reviewer_id = Column(Integer, ForeignKey("users.id"), nullable=False, index=True)
    feedback = Column(Text, nullable=False)
    issues_count = Column(Integer, default=0)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # Relationships
    project = relationship("Projects", back_populates="reviews")
    reviewer = relationship("Users")
    ai_results = relationship("AiResults", back_populates="review")

# ------------------ AI RESULTS ------------------
class AiResults(Base):
    __tablename__ = "ai_results"

    id = Column(Integer, primary_key=True, index=True)
    review_id = Column(Integer, ForeignKey("reviews.id"), nullable=False)
    output = Column(Text, nullable=False)
    model_version = Column(String(50), nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    # Relationships
    review = relationship("Reviews", back_populates="ai_results")
